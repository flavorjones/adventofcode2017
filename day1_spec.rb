#! /usr/bin/env ruby

require 'rspec'

class CaptchaSolver
  attr :input

  def initialize input
    @input = input.split('').map { |char| char.to_i }
  end

  def solution1
    solve_with do |index|
      next_index_for index
    end
  end

  def solution2
    solve_with do |index|
      jump_index_for index
    end
  end

  private

  def solve_with &index_func
    sum = 0

    input.each_with_index do |number, place|
      following = input[index_func.call(place)]
      if number == following
        sum = sum + number
      end
    end

    sum
  end

  def next_index_for index
    next_index = index + 1
    next_index % input.length
  end

  def jump_index_for index
    jump_length = input.length / 2
    jump_index = index + jump_length

    jump_index % input.length
  end
end

describe CaptchaSolver do
  it { expect(CaptchaSolver.new("1122").solution1).to eq(3) }
  it { expect(CaptchaSolver.new("1111").solution1).to eq(4) }
  it { expect(CaptchaSolver.new("1234").solution1).to eq(0) }
  it { expect(CaptchaSolver.new("91212129").solution1).to eq(9) }

  it { expect(CaptchaSolver.new("1212").solution2).to eq(6) }
  it { expect(CaptchaSolver.new("1221").solution2).to eq(0) }
  it { expect(CaptchaSolver.new("123425").solution2).to eq(4) }
  it { expect(CaptchaSolver.new("123123").solution2).to eq(12) }
  it { expect(CaptchaSolver.new("12131415").solution2).to eq(4) }
end

d1p1_input = "237369991482346124663395286354672985457326865748533412179778188397835279584149971999798512279429268727171755461418974558538246429986747532417846157526523238931351898548279549456694488433438982744782258279173323381571985454236569393975735715331438256795579514159946537868358735936832487422938678194757687698143224139243151222475131337135843793611742383267186158665726927967655583875485515512626142935357421852953775733748941926983377725386196187486131337458574829848723711355929684625223564489485597564768317432893836629255273452776232319265422533449549956244791565573727762687439221862632722277129613329167189874939414298584616496839223239197277563641853746193232543222813298195169345186499866147586559781523834595683496151581546829112745533347796213673814995849156321674379644323159259131925444961296821167483628812395391533572555624159939279125341335147234653572977345582135728994395631685618135563662689854691976843435785879952751266627645653981281891643823717528757341136747881518611439246877373935758151119185587921332175189332436522732144278613486716525897262879287772969529445511736924962777262394961547579248731343245241963914775991292177151554446695134653596633433171866618541957233463548142173235821168156636824233487983766612338498874251672993917446366865832618475491341253973267556113323245113845148121546526396995991171739837147479978645166417988918289287844384513974369397974378819848552153961651881528134624869454563488858625261356763562723261767873542683796675797124322382732437235544965647934514871672522777378931524994784845817584793564974285139867972185887185987353468488155283698464226415951583138352839943621294117262483559867661596299753986347244786339543174594266422815794658477629829383461829261994591318851587963554829459353892825847978971823347219468516784857348649693185172199398234123745415271222891161175788713733444497592853221743138324235934216658323717267715318744537689459113188549896737581637879552568829548365738314593851221113932919767844137362623398623853789938824592"

d1p1_answer = CaptchaSolver.new(d1p1_input).solution1
puts "d1p1: #{d1p1_answer}"

d1p2_answer = CaptchaSolver.new(d1p1_input).solution2
puts "d1p2: #{d1p2_answer}"
